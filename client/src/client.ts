import {createConnectTransport} from "@bufbuild/connect-web";
import {createPromiseClient} from "@bufbuild/connect";
import {MainService} from "./pbgen/v1/main_connect";

const fetchConfig = async () => {
  try {
    const response = await fetch(`${window.location.origin}/ap/config.json`);
    if (response.ok) {
      const json = await response.json();
      return json.baseUrl;
    } else {
      return window.location.origin;
    }
  } catch (error) {
    return window.location.origin;
  }
};

export const baseUrl = process.env.REACT_APP_API_ENDPOINT ?? await fetchConfig()

console.log("baseUrl:", baseUrl);

const transport = createConnectTransport({
  baseUrl,
  credentials: "include",
});

export const client = createPromiseClient(MainService, transport);
