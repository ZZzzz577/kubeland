import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router";
import routers from "./routes";
import { useMemo } from "react";

function App() {
  const rs = useMemo(() => {
    return createBrowserRouter(routers);
  }, []);
  return <RouterProvider router={rs} />;
}

export default App;
