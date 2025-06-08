function subscribe(eventName: any, listener: (data: CustomEvent) => any) {
    document.addEventListener(eventName, listener);
}

function unsubscribe(eventName: any, listener: (data: CustomEvent) => any) {
    document.removeEventListener(eventName, listener);
}

function dispatch(eventName: any, data: { code: string | number, message: string }) {
    const event = new CustomEvent(eventName, { detail: data });
    document.dispatchEvent(event);
}

export { dispatch, subscribe, unsubscribe };