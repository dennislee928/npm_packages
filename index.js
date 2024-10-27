import { get } from "axios"; // 引入 axios

function helpMessage() {
  const message = "Hey step-brother, help me, I am stucked!";
  console.log(message);
  return message;
}

// 新增一個函數來獲取隨機pornhub影片
async function fetchRandomVideo() {
  try {
    const response = await get("https://lust.scathach.id/pornhub/random");
    if (response.data.success) {
      return response.data.data; // 返回影片數據
    } else {
      throw new Error("Failed to fetch video");
    }
  } catch (error) {
    console.error(error);
    return null; // 錯誤處理
  }
}

// 新增一個函數來獲取隨機 Redtube 影片
async function fetchRandomRedtubeVideo() {
  try {
    const response = await axios.get("https://lust.scathach.id/redtube/random");
    if (response.data.success) {
      return response.data.data; // 返回影片數據
    } else {
      throw new Error("Failed to fetch video");
    }
  } catch (error) {
    console.error(error);
    return null; // 錯誤處理
  }
}

export default { helpMessage, fetchRandomVideo, fetchRandomRedtubeVideo }; // 將新函數導出
