// rxjs
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

// membership
import { LoginRequest, RegisterRequest } from "../components/shared/types";
import { ENV } from "../env";
import { Jwt } from "../components/shared/types/auth.interface";
import { HTTPService } from "./http.service";

export class AuthService extends HTTPService {
  private readonly authUrlSegment: string = ENV.api + "/auth";

  registerUser(request: RegisterRequest): Observable<void> {
    return this.post<void>(this.authUrlSegment + "/register", request);
  }

  login(request: LoginRequest): Observable<Jwt> {
    return this.post<Jwt>(this.authUrlSegment + "/login", request);
  }

  logout(): Observable<null> {
    return this.post<null>(this.authUrlSegment + "/logout").pipe(
      map(() => null)
    );
  }
}
