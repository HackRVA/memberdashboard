export type AuthResponse = {
  token: string;
};

export type AuthUser = {
  token: string;
  isLogin: boolean;
  email: string;
};

export type RegisterRequest = {
  email: string;
  password: string;
};

export type LoginRequest = RegisterRequest;
