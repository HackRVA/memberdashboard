export interface RegisterRequest {
  password: string;
  email: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface Jwt {
  token: string;
}

export interface UserProfile {
  email: string;
}
