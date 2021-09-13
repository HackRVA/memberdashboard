import { HTTPRequest } from './../types/api/http-request';
// rxjs
import { Observable, throwError } from 'rxjs';
import { fromFetch } from 'rxjs/fetch';
import { switchMap, catchError } from 'rxjs/operators';

export class HTTPService {
  get<T>(
    endpoint: string,
    responseType: 'json' | 'blob' = 'json'
  ): Observable<T> {
    return this.request<T>({
      methodName: 'GET',
      url: endpoint,
      responseType: responseType,
      options: {
        authType: 'Bearer',
      },
    } as HTTPRequest);
  }

  post<T>(
    endpoint: string,
    body?: any,
    authType: 'Basic' | 'Bearer' = 'Bearer'
  ): Observable<T> {
    return this.request<T>({
      methodName: 'POST',
      url: endpoint,
      responseType: 'json',
      options: {
        body: body,
        authType: authType,
      },
    } as HTTPRequest);
  }

  delete<T>(endpoint: string, body?: any): Observable<T> {
    return this.request<T>({
      methodName: 'DELETE',
      url: endpoint,
      responseType: 'json',
      options: {
        body: body,
        authType: 'Bearer',
      },
    } as HTTPRequest);
  }

  put<T>(endpoint: string, body?: any): Observable<T> {
    return this.request<T>({
      methodName: 'PUT',
      url: endpoint,
      responseType: 'json',
      options: {
        body: body,
        authType: 'Bearer',
      },
    } as HTTPRequest);
  }

  private request<R>(httpRequest: HTTPRequest): Observable<R> {
    return fromFetch(httpRequest.url, {
      method: httpRequest.methodName,
      headers: {
        Authorization: this.authHeader(httpRequest.options),
      },
      body: JSON.stringify(httpRequest.options.body),
    }).pipe(
      switchMap((response: Response) => {
        if (response.ok) {
          return this.formatResponse(response, httpRequest.responseType);
        }

        return throwError({
          error: true,
          message: `Error ${response.status}`,
        });
      }),
      catchError((err: Error) => {
        // Network or other error, handle appropriately
        return throwError({ error: true, message: err.message });
      })
    );
  }

  private formatResponse(
    response: Response,
    responseType: 'json' | 'blob'
  ): Promise<any> {
    switch (responseType) {
      case 'json':
        return response.json();
      case 'blob':
        return response.blob();
      default:
        throw 'Response type is not supported';
    }
  }

  private authHeader(options: {
    body?: any;
    authType: 'Basic' | 'Bearer';
  }): string {
    switch (options.authType) {
      case 'Basic':
        return (
          'Basic ' + btoa(options.body.email + ':' + options.body.password)
        );
      case 'Bearer':
        return 'Bearer ' + localStorage.getItem('jwt');
      default:
        return '';
    }
  }
}
