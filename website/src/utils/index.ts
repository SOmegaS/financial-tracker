interface ErrorWithMessage {
    message: string;
}

function isErrorWithMessage(error: unknown): error is ErrorWithMessage {
    return typeof error === 'object' && error !== null && 'message' in error;
}

export function getErrorMessage(error: unknown): string {
    if (isErrorWithMessage(error)) {
        return error.message;
    }
    return 'Что-то пошло не так';
}

export function getUserIdFromToken(token: string): string {
    const base64Url = token.split('.')[1];

    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');

    const payload = JSON.parse(atob(base64));
    return payload.user_id || null;
}
