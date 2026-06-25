export type UserRole = "user" | "admin";

export interface PublicUser {
  id: string;
  name: string;
  email: string | null;
  phone: string | null;
  role: UserRole;
  contributionScore: number;
  createdAt: string;
}

export interface AuthResult {
  user: PublicUser;
  token: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  data: AuthResult;
}

export interface ApiErrorResponse {
  success: false;
  message: string;
  error: string;
}

export interface RegisterPayload {
  name: string;
  email?: string;
  phone?: string;
  password: string;
}

export interface LoginPayload {
  email?: string;
  phone?: string;
  password: string;
}
