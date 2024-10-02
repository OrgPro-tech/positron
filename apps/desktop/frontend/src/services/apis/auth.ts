import { APIResponse, fetchClient, handleFetchError } from ".";

type LoginUser = {
  email: string;
  id: number;
  name: string;
};

export type User = {
  email: string;
  id: string;
  is_active: boolean;
  name: string;
  role: string;
};

export const emailAndPasswordLogin = async (
  email: string,
  password: string,
) => {
  try {
    const res = await fetchClient<APIResponse<LoginUser>>("login", {
      method: "post",
      body: {
        email: email,
        password: password,
      },
    });
    return res.data;
  } catch (e) {
    return handleFetchError(e);
  }
};
