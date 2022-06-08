export const DefaultSort = (asc: boolean, left: unknown, right: unknown): number => {
    if (asc) {
        if (left > right) {
            return 1;
        } else if (left === right) {
            return 0;
        }
        return -1;
    }
    if (left > right) {
        return -1;
    } else if (left === right) {
        return 0;
    }
    return 1;
};

export const DateSort = (asc: boolean, left: string, right: string): number => {
    const leftDate = new Date(left).valueOf();
    const rightDate = new Date(right).valueOf();

    return DefaultSort(asc, leftDate, rightDate);
};
