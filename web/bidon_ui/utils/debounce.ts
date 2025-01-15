export function debounce<T extends (...args: any[]) => void>(
  func: T,
  delay: number,
) {
  let timeout: ReturnType<typeof setTimeout>;

  return function (...args: Parameters<T>) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), delay);
  };
}
