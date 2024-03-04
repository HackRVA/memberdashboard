export function withSafeHandler<T>(errorId: string, callback: () => T): T {
  try {
    return callback();
  } catch (error) {
    console.error('Failed operation', errorId, error);
    return null;
  }
}
