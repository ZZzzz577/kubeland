export const getActivePath = (prefix: string, path: string): string => {
    return path.replace(prefix, "").replace(/^\/+|\/+$/g, "");
};

export default getActivePath;