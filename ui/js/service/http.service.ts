// rxjs
import { Observable, throwError } from "rxjs";
import { fromFetch } from "rxjs/fetch";
import { switchMap, catchError } from "rxjs/operators";

export class HTTPService {
  get(
    endpoint: string
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint, {
      method: "GET",
      headers: {
        Authorization: "Bearer " + localStorage.getItem("jwt"),
      },
    }).pipe(
      switchMap((response: Response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return throwError({
            error: true,
            message: `Error ${response.status}`,
          });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return throwError({ error: true, message: err.message });
      })
    );
  }
  post(
    endpoint: string,
    options?: any
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint, {
      method: "POST",
      headers: {
        Authorization: "Bearer " + localStorage.getItem("jwt"),
      },
      body: JSON.stringify(options),
    }).pipe(
      switchMap((response: Response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return throwError({
            error: true,
            message: `Error ${response.status}`,
          });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return throwError({ error: true, message: err.message });
      })
    );
  }
  delete(
    endpoint: string,
    options?: any
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint, {
      method: "DELETE",
      headers: {
        Authorization: "Bearer " + localStorage.getItem("jwt"),
      },
      body: JSON.stringify(options),
    }).pipe(
      switchMap((response: Response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return throwError({
            error: true,
            message: `Error ${response.status}`,
          });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return throwError({ error: true, message: err.message });
      })
    );
  }
  put(
    endpoint: string,
    options?: any
  ): Observable<Response | any | { error: boolean; message: any }> {
    return fromFetch(endpoint, {
      method: "PUT",
      headers: {
        Authorization: "Bearer " + localStorage.getItem("jwt"),
      },
      body: JSON.stringify(options),
    }).pipe(
      switchMap((response: Response) => {
        if (response.ok) {
          // OK return data
          return response.json();
        } else {
          // Server is returning a status requiring the client to try something else.
          return throwError({
            error: true,
            message: `Error ${response.status}`,
          });
        }
      }),
      catchError((err) => {
        // Network or other error, handle appropriately
        console.error(err);
        return throwError({ error: true, message: err.message });
      })
    );
  }
}
