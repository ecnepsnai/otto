import { Notification } from '../components/Notification';

/**
 * Class for interacting with the API
 */
export class API {
    private static do(url: string, options: RequestInit, noCheckError?: boolean): Promise<unknown> {
        return new Promise((resolve, reject) => {
            fetch(url, options).then(response => {
                response.json().then(results => {
                    if (response.status === 403) {
                        location.href = '/login?unauthorized&redirect=' + location.pathname;
                        return;
                    }

                    if (noCheckError) {
                        resolve(results.data);
                        return;
                    }

                    if (!response.ok) {
                        let message = 'Internal Server Error';
                        if (results.error && results.error.message) {
                            message = results.error.message;
                        }
                        console.error('API error caught', results);
                        Notification.error('An Error Occurred: ' + message);
                        reject(results);
                    } else {
                        resolve(results.data);
                    }
                }, e => {
                    if (!noCheckError) {
                        console.error('API error caught', e);
                        Notification.error('An Error Occurred: Internal Server Error');
                    }
                    reject(e);
                });
            }, e => {
                if (!noCheckError) {
                    console.error('API error caught', e);
                    Notification.error('An Error Occurred: Internal Server Error');
                }
                reject(e);
            });
        });
    }

    /**
     * Perform a HTTP GET request to the specified URL
     * @param url the URL to request
     * @returns The JSON object of the results
     */
    public static GET(url: string): Promise<unknown> {
        return this.do(url, { method: 'GET' });
    }

    /**
     * Perform a HTTP GET request to the specified URL but do
     * not handle any errors
     * @param url the URL to request
     * @returns The JSON object of the results
     */
    public static UnsafeGET(url: string): Promise<unknown> {
        return this.do(url, { method: 'GET' }, true);
    }

    /**
     * Perform a HTTP POST request to the specified URL with the given body
     * @param url the URL to request
     * @param data Body data to be encoded as JSON
     * @returns The JSON object of the results
     */
    public static async POST(url: string, data: unknown): Promise<unknown> {
        return this.do(url, { method: 'POST', body: JSON.stringify(data) });
    }

    /**
     * Perform a HTTP PUT request to the specified URL with the given body
     * @param url the URL to request
     * @param data Body data to be encoded as JSON
     * @returns The JSON object of the results
     */
    public static async PUT(url: string, data: unknown): Promise<unknown> {
        return this.do(url, { method: 'PUT', body: JSON.stringify(data) });
    }

    /**
     * Perform a HTTP PATCH request to the specified URL with the given body
     * @param url the URL to request
     * @param data Body data to be encoded as JSON
     * @returns The JSON object of the results
     */
    public static async PATCH(url: string, data: unknown): Promise<unknown> {
        return this.do(url, { method: 'PATCH', body: JSON.stringify(data) });
    }

    /**
     * Perform a HTTP DELETE request to the specified URL
     * @param url the URL to request
     * @returns The JSON object of the results
     */
    public static async DELETE(url: string): Promise<unknown> {
        return this.do(url, { method: 'DELETE' });
    }

    private static async UploadFile(method: 'POST' | 'PUT', url: string, file: File, data: { [key: string]: string; }): Promise<unknown> {
        const fd = new FormData();
        if (file) {
            fd.append('file', file);
        }
        Object.keys(data).forEach(key => {
            fd.append(key, data[key]);
        });

        return this.do(url, {
            method: method,
            body: fd
        });
    }

    /**
     * Perform a HTTP PUT Multipart Upload to the specified URL
     * @param url the URL to request
     * @param file The file data
     * @param data Additional parameters to add to the form data
     * @returns The JSON object of the results
     */
    public static async PUTFile(url: string, file: File, data: { [key: string]: string; }): Promise<unknown> {
        return API.UploadFile('PUT', url, file, data);
    }

    /**
     * Perform a HTTP POST Multipart Upload to the specified URL
     * @param url the URL to request
     * @param file The file data
     * @param data Additional parameters to add to the form data
     * @returns The JSON object of the results
     */
    public static async POSTFile(url: string, file: File, data: { [key: string]: string; }): Promise<unknown> {
        return API.UploadFile('POST', url, file, data);
    }
}
