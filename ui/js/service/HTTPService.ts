import { of, Observable } from "rxjs";
import { fromFetch } from "rxjs/fetch";
import { switchMap, catchError } from "rxjs/operators";

export class HTTPService {
  get(
    endpoint: string
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint).pipe(
      switchMap((response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return of({ error: true, message: `Error ${response.status}` });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return of({ error: true, message: err.message });
      })
    );
  }
  post(
    endpoint: string,
    options?: any
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint, {
      method: "POST",
      body: JSON.stringify(options),
    }).pipe(
      switchMap((response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return of({ error: true, message: `Error ${response.status}` });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return of({ error: true, message: err.message });
      })
    );
  }
}
