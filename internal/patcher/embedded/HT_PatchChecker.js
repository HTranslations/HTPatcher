/*:
 * @target MV MZ
 * @plugindesc [v1.1] HTranslations — patch update checker + auto word wrap.
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
 * @param checkForUpdates
 * @text Check For Updates
 * @desc Check for translation patch updates on startup.
 * @type boolean
 * @default true
 *
 * @param autoWrap
 * @text Auto Word Wrap
 * @desc Automatically wrap message text to fit the window width.
 * @type boolean
 * @default false
 *
 * @help
 * ============================================================================
 * HTranslations Plugin
 * ============================================================================
 *
 * This plugin provides two features:
 *
 * 1. Patch Update Checker — checks for translation patch updates when the
 *    game starts. If a newer version is available, it displays the patch
 *    notes and offers a link to download the update.
 *
 * 2. Auto Word Wrap — automatically wraps message text to fit the message
 *    window width at display time, accounting for face graphics. This
 *    eliminates the need for a fixed wrap width at export time.
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
  var CHECK_UPDATES = String(params["checkForUpdates"] || "true") === "true";
  var AUTO_WRAP = String(params["autoWrap"] || "false") === "true";

  // -------------------------------------------------------------------------
  // Auto Word Wrap — wraps message text to fit the window at display time
  // -------------------------------------------------------------------------
  if (AUTO_WRAP) {

    var _HT_WM_startMessage = Window_Message.prototype.startMessage;
    Window_Message.prototype.startMessage = function () {
      _HT_WM_startMessage.call(this);
      if (this._textState && this._textState.text && this.contents) {
        this._textState.text = this._htAutoWrap(this._textState.text);
      }
    };

    Window_Message.prototype._htAutoWrap = function (text) {
      this.resetFontSettings();
      var startX = this.newLineX(this._textState);
      var maxWidth = this.contentsWidth() - startX;
      if (maxWidth <= 0) return text;

      var lines = text.split("\n");
      var wrapped = [];
      for (var i = 0; i < lines.length; i++) {
        wrapped.push(this._htWrapLine(lines[i], maxWidth));
      }
      return wrapped.join("\n");
    };

    /**
     * Read an escape sequence starting at text[i] (\x1b).
     * Handles multi-letter codes (e.g. FS, PX) and special chars ({, }, $).
     * Returns { end: next index, width: visible pixel width }.
     */
    Window_Message.prototype._htReadEscape = function (text, i) {
      i++; // skip \x1b
      var codeStart = i;
      while (i < text.length) {
        var cc = text.charCodeAt(i);
        if ((cc >= 65 && cc <= 90) || (cc >= 97 && cc <= 122)) { i++; } else { break; }
      }
      var code = text.substring(codeStart, i).toUpperCase();
      if (i === codeStart && i < text.length) i++; // single special char
      if (i < text.length && text[i] === "[") {
        i++;
        while (i < text.length && text[i] !== "]") i++;
        if (i < text.length) i++;
      }
      var w = 0;
      if (code === "I") w = (ImageManager.iconWidth || 32) + 4;
      return { end: i, width: w };
    };

    /** Extract only visible characters from text (strip escape sequences). */
    Window_Message.prototype._htVisibleText = function (text) {
      var visible = "";
      var i = 0;
      while (i < text.length) {
        if (text[i] === "\x1b") {
          i = this._htReadEscape(text, i).end;
        } else {
          visible += text[i];
          i++;
        }
      }
      return visible;
    };

    /** Sum escape-based pixel width (icons) in text. */
    Window_Message.prototype._htEscapeWidth = function (text) {
      var w = 0;
      var i = 0;
      while (i < text.length) {
        if (text[i] === "\x1b") {
          var esc = this._htReadEscape(text, i);
          w += esc.width;
          i = esc.end;
        } else {
          i++;
        }
      }
      return w;
    };

    /**
     * Wrap a single line to fit within maxWidth pixels.
     * Uses word wrap with whole-string measurement for accurate kerning.
     */
    Window_Message.prototype._htWrapLine = function (line, maxWidth) {
      // Tokenize into words and spaces (escape sequences attach to adjacent words)
      var tokens = [];
      var i = 0;
      var word = "";
      while (i < line.length) {
        if (line[i] === "\x1b") {
          var esc = this._htReadEscape(line, i);
          word += line.substring(i, esc.end);
          i = esc.end;
        } else if (line[i] === " ") {
          if (word) { tokens.push(word); word = ""; }
          tokens.push(" ");
          i++;
        } else {
          word += line[i];
          i++;
        }
      }
      if (word) tokens.push(word);

      // Build wrapped output, tracking the current line's visible text
      // so we can measure the whole line at once (accurate kerning)
      var lines = [[]];
      var curVisible = "";
      var curEscW = 0;

      for (var t = 0; t < tokens.length; t++) {
        var tok = tokens[t];

        if (tok === " ") {
          var testWidth = this.contents.measureTextWidth(curVisible + " ") + curEscW;
          if (testWidth > maxWidth) {
            lines.push([]);
            curVisible = "";
            curEscW = 0;
          } else {
            lines[lines.length - 1].push(" ");
            curVisible += " ";
          }
          continue;
        }

        var tokVisible = this._htVisibleText(tok);
        var tokEscW = this._htEscapeWidth(tok);
        var testWidth2 = this.contents.measureTextWidth(curVisible + tokVisible) + curEscW + tokEscW;

        if (curVisible.length > 0 && testWidth2 > maxWidth) {
          // Trim trailing space from current line
          var curLine = lines[lines.length - 1];
          if (curLine.length > 0 && curLine[curLine.length - 1] === " ") {
            curLine.pop();
            curVisible = curVisible.substring(0, curVisible.length - 1);
          }
          lines.push([]);
          curVisible = "";
          curEscW = 0;
        }

        // Word alone is wider than maxWidth — character-wrap it
        if (this.contents.measureTextWidth(tokVisible) + tokEscW > maxWidth) {
          var ci = 0;
          while (ci < tok.length) {
            if (tok[ci] === "\x1b") {
              var cesc = this._htReadEscape(tok, ci);
              if (cesc.width > 0) {
                var cTestW = this.contents.measureTextWidth(curVisible) + curEscW + cesc.width;
                if (cTestW > maxWidth && curVisible.length > 0) {
                  lines.push([]);
                  curVisible = "";
                  curEscW = 0;
                }
                curEscW += cesc.width;
              }
              lines[lines.length - 1].push(tok.substring(ci, cesc.end));
              ci = cesc.end;
            } else {
              var cTestW2 = this.contents.measureTextWidth(curVisible + tok[ci]) + curEscW;
              if (cTestW2 > maxWidth && curVisible.length > 0) {
                lines.push([]);
                curVisible = "";
                curEscW = 0;
              }
              curVisible += tok[ci];
              lines[lines.length - 1].push(tok[ci]);
              ci++;
            }
          }
          continue;
        }

        lines[lines.length - 1].push(tok);
        curVisible += tokVisible;
        curEscW += tokEscW;
      }

      var output = [];
      for (var li = 0; li < lines.length; li++) {
        output.push(lines[li].join(""));
      }
      return output.join("\n");
    };
  }

  // -------------------------------------------------------------------------
  // Patch Checker — fetch and display update notifications
  // -------------------------------------------------------------------------
  if (!CHECK_UPDATES || !STORE_CODE) return;

  var _updateData = null;
  var _fetchDone = false;
  var _shown = false;

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
          if (CURRENT_VERSION === "0") {
            var allPatches = [];
            for (var i = 0; i < data.patches.length; i++) {
              var p = data.patches[i];
              allPatches.push({ version: p.version, releaseNotes: p.releaseNotes || [] });
            }
            _updateData = {
              latestVersion: latest.version,
              missedPatches: allPatches,
              slug: data.slug || "",
              manuallyPatched: true,
            };
          } else if (String(latest.version) !== String(CURRENT_VERSION)) {
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
              manuallyPatched: false,
            };
          }
        }
      } catch (e) {}
    }
    _fetchDone = true;
  };
  xhr.onerror = function () { _fetchDone = true; };
  xhr.send();

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

  var _Scene_Title_isBusy = Scene_Title.prototype.isBusy;
  Scene_Title.prototype.isBusy = function () {
    if (this._htWindow && this._htWindow.isOpen()) return true;
    return _Scene_Title_isBusy.call(this);
  };

  // -------------------------------------------------------------------------
  // Window_HTPatch — scrollable update notification window
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

    if (this._htData.manuallyPatched) {
      lines.push({ text: "Manually Patched", y: 0, align: "center", color: "#f5c542", size: 20 });
      lines.push({ text: "Current version could not be determined", y: 34, align: "center", color: "#aaaaaa", size: 14 });
    } else {
      lines.push({ text: "Patch Update Available!", y: 0, align: "center", color: "#f5c542", size: 20 });
      lines.push({ text: "v" + CURRENT_VERSION + "  \u2192  v" + this._htData.latestVersion, y: 34, align: "center", color: "#ffffff", size: 16 });
    }

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
      cy += 8;
    }

    this._textLines = lines;
    this._totalHeight = cy;
    var footerH = 56;
    var pad = typeof this.standardPadding === "function" ? this.standardPadding() : this.padding;
    var scrollableH = this.height - pad * 2 - footerH;
    this._scrollableH = scrollableH;
    this._maxScroll = Math.max(0, this._totalHeight - scrollableH);
  };

  Window_HTPatch.prototype.refresh = function () {
    this.createContents();
    if (!this._textLines) return;
    var innerW = this.contentsWidth();
    var scrollableH = this._scrollableH;

    for (var i = 0; i < this._textLines.length; i++) {
      var l = this._textLines[i];
      var dy = l.y - this._scrollY;
      if (dy + (l.size + 8) <= 0 || dy >= scrollableH) continue;
      this.contents.fontSize = l.size || 14;
      this.contents.textColor = l.color || "#ffffff";
      if (this.changeTextColor) this.changeTextColor(l.color || "#ffffff");
      this.contents.drawText(l.text, 0, dy, innerW, l.size + 8, l.align || "left");
    }

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
