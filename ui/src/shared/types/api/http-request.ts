export interface HTTPRequest {
  methodName: 'GET' | 'POST' | 'PUT' | 'DELETE';
  options: { authType: 'Basic' | 'Bearer'; body?: any };
  url: string;
  responseType: 'json' | 'blob';
}
