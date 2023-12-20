export default function useParseDate() {
  const leftPad = (num: number) => num.toString().padStart(2, "0");
  // Parse date to YYYY/MM/DD HH:mm:ss (Server format)
  return (
    date: string,
    options: { dateOnly: boolean } = { dateOnly: false }
  ) => {
    //to YYYY/MM/DD HH:mm:ss
    const dateObj = new Date(date);
    const year = dateObj.getFullYear();
    const month = dateObj.getMonth() + 1;
    const day = dateObj.getDate();
    const hours = dateObj.getHours();
    const minutes = dateObj.getMinutes();
    const seconds = dateObj.getSeconds();
    const formattedTime =
      `${year}/${leftPad(month)}/${leftPad(day)}` +
      (options.dateOnly
        ? ``
        : ` ${leftPad(hours)}:${leftPad(minutes)}:${leftPad(seconds)}`);

    return formattedTime;
  };
}
