export function getDlsiteImageUrl(code: string): string {
  if (!code || !code.startsWith("RJ")) return "";
  const numberStr = code.substring(2);
  const number = parseInt(numberStr);
  if (isNaN(number)) return "";
  const folderNumber = Math.ceil(number / 1000) * 1000;
  const folderCode = `RJ${folderNumber.toString().padStart(numberStr.length, "0")}`;
  return `https://img.dlsite.jp/modpub/images2/work/doujin/${folderCode}/${code}_img_main.webp`;
}

