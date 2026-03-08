/*:
 * @target MV MZ
 * @plugindesc [v1.0] HTranslations Patch Checker — checks for translation updates on startup.
 * @author HTranslations
 * @url https://htranslations.com
 *
 * @param storeCode
 * @text Store Code
 * @desc The store code for this game (e.g. RJ405647)
 * @type string
 * @default
 *
 * @param currentVersion
 * @text Current Patch Version
 * @desc The version of the currently installed patch.
 * @type string
 * @default 1
 *
 * @param apiBase
 * @text API Base URL
 * @desc Base URL of the HTranslations API.
 * @type string
 * @default https://htranslations.com
 *
 * @help
 * ============================================================================
 * HTranslations Patch Checker
 * ============================================================================
 *
 * This plugin checks for translation patch updates when the game starts.
 * If a newer version is available, it displays the patch notes and offers
 * a link to download the update.
 *
 * Configure the Store Code and Current Patch Version in the plugin parameters.
 * The Current Patch Version should be updated each time you release a new patch.
 *
 * Compatible with RPG Maker MV and MZ.
 * ============================================================================
 */

(function () {
  "use strict";

  var params = PluginManager.parameters("HT_PatchChecker");
  var STORE_CODE = String(params["storeCode"] || "");
  var CURRENT_VERSION = String(params["currentVersion"] || "1");
  var API_BASE = String(params["apiBase"] || "https://htranslations.com");

  if (!STORE_CODE) return;

  var _updateData = null;
  var _fetchDone = false;
  var _shown = false;

  // Fire fetch immediately on plugin load
  var url = API_BASE + "/api/patches/" + encodeURIComponent(STORE_CODE);
  var xhr = new XMLHttpRequest();
  xhr.open("GET", url);
  xhr.overrideMimeType("application/json");
  xhr.onload = function () {
    if (xhr.status === 200) {
      try {
        var data = JSON.parse(xhr.responseText);
        if (data.patches && data.patches.length > 0) {
          var latest = data.patches[0];
          if (String(latest.version) !== String(CURRENT_VERSION)) {
            // Collect all patches newer than the installed version
            var missed = [];
            for (var i = 0; i < data.patches.length; i++) {
              var p = data.patches[i];
              if (String(p.version) === String(CURRENT_VERSION)) break;
              missed.push({ version: p.version, releaseNotes: p.releaseNotes || [] });
            }
            _updateData = {
              latestVersion: latest.version,
              missedPatches: missed,
              slug: data.slug || "",
            };
          }
        }
      } catch (e) {}
    }
    _fetchDone = true;
  };
  xhr.onerror = function () { _fetchDone = true; };
  xhr.send();

  // Hook Scene_Title.update to show window once fetch is done
  var _Scene_Title_update = Scene_Title.prototype.update;
  Scene_Title.prototype.update = function () {
    _Scene_Title_update.call(this);
    if (_fetchDone && _updateData && !_shown) {
      _shown = true;
      var ww = Math.min(Graphics.boxWidth - 80, 560);
      var wh = Math.min(Graphics.boxHeight - 60, 400);
      var wx = (Graphics.boxWidth - ww) / 2;
      var wy = (Graphics.boxHeight - wh) / 2;
      this._htWindow = new Window_HTPatch(wx, wy, ww, wh, _updateData);
      this.addChild(this._htWindow);
    }
  };

  // Also block title command window input while our window is open
  var _Scene_Title_isBusy = Scene_Title.prototype.isBusy;
  Scene_Title.prototype.isBusy = function () {
    if (this._htWindow && this._htWindow.isOpen()) return true;
    return _Scene_Title_isBusy.call(this);
  };

  // -------------------------------------------------------------------------
  // Window_HTPatch — simple scrollable text window
  // -------------------------------------------------------------------------
  function Window_HTPatch(x, y, w, h, data) {
    this._htData = data;
    this._scrollY = 0;
    this._maxScroll = 0;
    this._textLines = null;
    this.initialize(x, y, w, h);
  }

  Window_HTPatch.prototype = Object.create(Window_Base.prototype);
  Window_HTPatch.prototype.constructor = Window_HTPatch;

  Window_HTPatch.prototype.initialize = function (x, y, w, h) {
    if (Window_Base.prototype.initialize.length === 1) {
      Window_Base.prototype.initialize.call(this, new Rectangle(x, y, w, h));
    } else {
      Window_Base.prototype.initialize.call(this, x, y, w, h);
    }
    this.backOpacity = 255;
    this._buildText();
    this.refresh();
    this.open();
    this.activate();
  };

  Window_HTPatch.prototype.standardBackOpacity = function () {
    return 255;
  };


  Window_HTPatch.prototype._buildText = function () {
    var lines = [];
    var missed = this._htData.missedPatches || [];

    lines.push({ text: "Patch Update Available!", y: 0, align: "center", color: "#f5c542", size: 20 });
    lines.push({ text: "v" + CURRENT_VERSION + "  \u2192  v" + this._htData.latestVersion, y: 34, align: "center", color: "#ffffff", size: 16 });

    var cy = 68;
    var lh = 24;

    for (var p = 0; p < missed.length; p++) {
      var patch = missed[p];
      lines.push({ text: "v" + patch.version, y: cy, color: "#f5c542", size: 15 });
      cy += lh;
      var notes = patch.releaseNotes;
      if (Array.isArray(notes)) {
        for (var i = 0; i < notes.length; i++) {
          lines.push({ text: "\u2022 " + notes[i], y: cy, color: "#dddddd", size: 14 });
          cy += lh;
        }
      }
      cy += 8; // gap between patches
    }

    this._textLines = lines;
    this._totalHeight = cy;
    // Reserve space for the fixed footer (2 lines)
    var footerH = 56;
    var scrollableH = this.height - this.standardPadding() * 2 - footerH;
    this._scrollableH = scrollableH;
    this._maxScroll = Math.max(0, this._totalHeight - scrollableH);
  };

  Window_HTPatch.prototype.refresh = function () {
    this.createContents();
    if (!this._textLines) return;
    var innerW = this.contentsWidth();
    var scrollableH = this._scrollableH;

    // Draw scrollable patch notes (clipped to scrollable area)
    for (var i = 0; i < this._textLines.length; i++) {
      var l = this._textLines[i];
      var dy = l.y - this._scrollY;
      if (dy + (l.size + 8) <= 0 || dy >= scrollableH) continue;
      this.contents.fontSize = l.size || 14;
      this.contents.textColor = l.color || "#ffffff";
      if (this.changeTextColor) this.changeTextColor(l.color || "#ffffff");
      this.contents.drawText(l.text, 0, dy, innerW, l.size + 8, l.align || "left");
    }

    // Draw fixed footer at bottom
    var footerY = this.contentsHeight() - 50;
    this.contents.fontSize = 13;
    this.contents.textColor = "#aaaaaa";
    if (this.changeTextColor) this.changeTextColor("#aaaaaa");
    this.contents.drawText("Visit htranslations.com to download the update or repatch using HTPatcher.", 0, footerY, innerW, 20, "center");
    this.contents.fontSize = 13;
    this.contents.textColor = "#666666";
    if (this.changeTextColor) this.changeTextColor("#666666");
    this.contents.drawText("[ OK / Enter to close ]", 0, footerY + 24, innerW, 20, "center");
  };

  // Override createContents to make a tall bitmap for scrolling
  Window_HTPatch.prototype.createContents = function () {
    var w = this.contentsWidth();
    var h = Math.max(this._totalHeight || 400, this.contentsHeight());
    this.contents = new Bitmap(w, h);
    this.resetFontSettings();
  };

  Window_HTPatch.prototype.update = function () {
    Window_Base.prototype.update.call(this);
    if (!this.isOpen()) return;

    var scrolled = false;
    if (Input.isPressed("down")) {
      this._scrollY = Math.min(this._scrollY + 3, this._maxScroll);
      scrolled = true;
    }
    if (Input.isPressed("up")) {
      this._scrollY = Math.max(this._scrollY - 3, 0);
      scrolled = true;
    }
    if (typeof TouchInput !== "undefined") {
      if (TouchInput.wheelY > 0) {
        this._scrollY = Math.min(this._scrollY + 16, this._maxScroll);
        scrolled = true;
      } else if (TouchInput.wheelY < 0) {
        this._scrollY = Math.max(this._scrollY - 16, 0);
        scrolled = true;
      }
    }
    if (scrolled) this.refresh();

    if (Input.isTriggered("ok") || Input.isTriggered("cancel")) {
      SoundManager.playOk();
      this.close();
    }
  };
})();
