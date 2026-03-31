const API_BASE = "http://localhost:8888/api";

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  const json: ApiResponse<T> = await res.json();
  if (json.code !== 0) {
    throw new Error(json.message);
  }
  return json.data;
}

export const api = {
  sourceItems: {
    import: (data: { platform: string; sourceUrl: string }) =>
      request("/source-items/import", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/source-items?${new URLSearchParams(params)}`
      ),
    updateStatus: (id: number, status: string) =>
      request(`/source-items/status?id=${id}`, {
        method: "PUT",
        body: JSON.stringify({ status }),
      }),
    addTag: (id: number, tag: string) =>
      request(`/source-items/tag?id=${id}`, {
        method: "POST",
        body: JSON.stringify({ tag }),
      }),
  },
  products: {
    createFromSource: (sourceItemId: number) =>
      request("/products/create-from-source", {
        method: "POST",
        body: JSON.stringify({ sourceItemId }),
      }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/products?${new URLSearchParams(params)}`
      ),
    get: (id: number) => request(`/products/detail?id=${id}`),
    update: (id: number, data: any) =>
      request(`/products?id=${id}`, { method: "PUT", body: JSON.stringify(data) }),
    markReady: (id: number) =>
      request(`/products/ready?id=${id}`, { method: "PUT" }),
  },
  publishTasks: {
    create: (data: {
      productId: number;
      targetPlatform: string;
      categoryId: string;
      freightTemplate: string;
    }) => request("/publish-tasks", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/publish-tasks?${new URLSearchParams(params)}`
      ),
  },
};
