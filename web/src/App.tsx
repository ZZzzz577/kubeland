import { createBrowserRouter, RouterProvider } from "react-router";
import routes from "@/routes";
import ThemeProvider from "@/components/theme/provider";
import LocalesProvider from "@/components/locales/provider";
import { App as AntApp } from "antd";

function App() {
    const router = createBrowserRouter(routes);
    return (
        <AntApp>
            <ThemeProvider>
                <LocalesProvider>
                    <RouterProvider router={router} />
                </LocalesProvider>
            </ThemeProvider>
        </AntApp>
    );
}

export default App;
