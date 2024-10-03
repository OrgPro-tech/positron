import APIError from "@/exceptions/api-error";
import { BetterFetchError, createFetch } from "@better-fetch/fetch";

export const fetchClient = createFetch({
  baseURL: "http://localhost:8080/v1/api",
  throw: true,
});

export type APIResponse<TData> = {
  status: number;
  message: string;
  data: TData;
};

export type APIErrorResponse = {
  status: number;
  error: {
    message: string;
  };
};

export const handleFetchError = (e: unknown) => {
  if (e instanceof BetterFetchError) {
    const {
      error: { message },
    } = e.error as APIErrorResponse;
    throw new APIError(message, e.status);
  }
  throw e;
};
