import { fetchClient, handleFetchError } from ".";

type LoginUser = {
  access_token: string;
  refresh_token: string;
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
    const res = await fetchClient<LoginUser>("login", {
      method: "post",
      body: {
        username: email,
        password: password,
      },
    });
    return res;
  } catch (e) {
    return handleFetchError(e);
  }
};
