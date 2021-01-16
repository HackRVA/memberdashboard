import { Observable } from "rxjs";
import { HTTPService } from "./HTTPService";

export class UserService extends HTTPService {
  login(
    loginRequest: UserService.LoginRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/signin", loginRequest);
  }

  logout(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/edge/logout");
  }

  getUser(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/user");
  }

  registerUser(
    registerRequest: UserService.RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/register", registerRequest);
  }
}

export namespace UserService {
  export interface RegisterRequest {
    username: string;
    password: string;
    email: string;
  }

  export interface LoginRequest {
    username: string;
    password: string;
    updateCallback?: Function;
  }

  export interface UserProfile {
    username: string;
    email: string;
  }
}
