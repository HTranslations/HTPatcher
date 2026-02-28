const CDN_URL = "https://files.htranslations.com/htranslations";

export function getDlsiteImageUrl(code: string): string {
  if (!code || !code.startsWith("RJ")) return "";
  const numberStr = code.substring(2);
  const number = parseInt(numberStr);
  if (isNaN(number)) return "";
  const folderNumber = Math.ceil(number / 1000) * 1000;
  const folderCode = `RJ${folderNumber.toString().padStart(numberStr.length, "0")}`;
  return `https://img.dlsite.jp/modpub/images2/work/doujin/${folderCode}/${code}_img_main.webp`;
}

export function getDmmImageUrl(code: string): string {
  if (!code || !code.startsWith("d_")) return "";
  return `https://doujin-assets.dmm.co.jp/digital/game/${code}/${code}pl.jpg`;
}

export function getImageUrlFromCode(code: string): string {
  if (!code) return "";
  const upper = code.toUpperCase();
  if (upper.startsWith("RJ")) return getDlsiteImageUrl(upper);
  if (code.toLowerCase().startsWith("d_")) return getDmmImageUrl(code.toLowerCase());
  return "";
}

export function getThumbnailUrl(
  store: string,
  storeCode: string | undefined,
  thumbnailId: string | undefined,
  thumbnailFileName: string | undefined,
): string {
  // For DLsite, use the DLsite CDN directly
  if (store === "dlsite" && storeCode) {
    return getDlsiteImageUrl(storeCode.toUpperCase());
  }
  // For DMM, use the DMM CDN directly
  if (store === "dmm" && storeCode) {
    return getDmmImageUrl(storeCode.toLowerCase());
  }
  // For other stores, use the htranslations CDN
  if (thumbnailId && thumbnailFileName) {
    return `${CDN_URL}/public/${thumbnailId}/${thumbnailFileName}`;
  }
  return "";
}

export function getStoreLinkLabel(store: string): string {
  switch (store) {
    case "dlsite": return "View on DLsite";
    case "dmm": return "View on DMM";
    default: return "View Store Page";
  }
}
