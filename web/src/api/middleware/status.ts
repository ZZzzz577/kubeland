import type { Middleware, ResponseContext } from "@/generated";

export const StatusMiddleware: Middleware = {
    async post(context: ResponseContext): Promise<Response | void> {
        const status = context.response.status;
        if (status > 300) {
            throw await context.response.json();
        }
        return context.response;
    },
};
