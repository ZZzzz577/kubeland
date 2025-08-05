import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router";
import routers from "@/routes";
import ThemeProvider from "@/components/theme/provider";
import LocalesProvider from "@/components/locales/provider";

function App() {
  return (
    <ThemeProvider>
      <LocalesProvider>
        <RouterProvider router={createBrowserRouter(routers)} />
      </LocalesProvider>
    </ThemeProvider>
  );
}

export default App;
