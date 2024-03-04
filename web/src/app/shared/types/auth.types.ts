export type AuthResponse = {
  token: string;
};

export type AuthUser = {
  token: string;
  isLogin: boolean;
  email: string;
  isAdmin: boolean;
};

export type RegisterRequest = {
  email: string;
  password: string;
};

export type LoginRequest = RegisterRequest;
