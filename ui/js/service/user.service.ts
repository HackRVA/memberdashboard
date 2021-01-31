import { Observable } from "rxjs";
import { HTTPService } from "./http.service";

export class UserService extends HTTPService {
  login(
    request: UserService.LoginRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/login", request);
  }

  logout(): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/logout");
  }

  getUser(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/user");
  }

  registerUser(
    request: UserService.RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/register", request);
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

  export interface Jwt {
    token: string;
  }

  export interface UserProfile {
    username: string;
    email: string;
  }
}
