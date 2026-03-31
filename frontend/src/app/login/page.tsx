"use client";

import { Button } from "@/components/ui/button";

export default function LoginPage() {
  const handleGoogleLogin = () => {
    window.location.href = "http://localhost:8880/api/auth/google/login";
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-sm w-full space-y-6 p-8 bg-white rounded-xl shadow-md">
        <div className="text-center">
          <h1 className="text-2xl font-bold">代发工具</h1>
          <p className="text-gray-500 mt-2">电商一键代发运营工具</p>
        </div>
        <Button onClick={handleGoogleLogin} className="w-full" size="lg">
          Sign in with Google
        </Button>
      </div>
    </div>
  );
}
