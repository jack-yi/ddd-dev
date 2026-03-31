import { getToken, clearToken } from "./auth";

const API_BASE = "http://localhost:8888/api";
const USER_CENTER_API = "http://localhost:8880/api";

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

async function request<T>(base: string, path: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}${path}`, { headers, ...options });
  if (res.status === 401) {
    clearToken();
    if (typeof window !== "undefined") window.location.href = "/login";
    throw new Error("unauthorized");
  }
  const json: ApiResponse<T> = await res.json();
  if (json.code !== 0) {
    throw new Error(json.message);
  }
  return json.data;
}

export const api = {
  sourceItems: {
    import: (data: { platform: string; sourceUrl: string }) =>
      request(API_BASE, "/source-items/import", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/source-items?${new URLSearchParams(params)}`),
    updateStatus: (id: number, status: string) =>
      request(API_BASE, `/source-items/status?id=${id}`, { method: "PUT", body: JSON.stringify({ status }) }),
    addTag: (id: number, tag: string) =>
      request(API_BASE, `/source-items/tag?id=${id}`, { method: "POST", body: JSON.stringify({ tag }) }),
  },
  products: {
    createFromSource: (sourceItemId: number) =>
      request(API_BASE, "/products/create-from-source", { method: "POST", body: JSON.stringify({ sourceItemId }) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/products?${new URLSearchParams(params)}`),
    get: (id: number) => request(API_BASE, `/products/detail?id=${id}`),
    update: (id: number, data: any) =>
      request(API_BASE, `/products?id=${id}`, { method: "PUT", body: JSON.stringify(data) }),
    markReady: (id: number) =>
      request(API_BASE, `/products/ready?id=${id}`, { method: "PUT" }),
  },
  publishTasks: {
    create: (data: { productId: number; targetPlatform: string; categoryId: string; freightTemplate: string }) =>
      request(API_BASE, "/publish-tasks", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/publish-tasks?${new URLSearchParams(params)}`),
  },
  auth: {
    me: () => request<any>(USER_CENTER_API, "/auth/me"),
    checkInit: () => request<{ needInit: boolean }>(USER_CENTER_API, "/init/check"),
    initSuperAdmin: () => request(USER_CENTER_API, "/init/super-admin", { method: "POST" }),
  },
  users: {
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(USER_CENTER_API, `/users?${new URLSearchParams(params)}`),
    updateStatus: (id: number, status: string) =>
      request(USER_CENTER_API, `/users/status?id=${id}`, { method: "PUT", body: JSON.stringify({ status }) }),
    assignRole: (id: number, roleName: string) =>
      request(USER_CENTER_API, `/users/role?id=${id}`, { method: "PUT", body: JSON.stringify({ roleName }) }),
  },
};
