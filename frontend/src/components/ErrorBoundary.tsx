import * as React from 'react';

interface ErrorBoundaryProps {
    children?: React.ReactNode
}
interface ErrorBoundaryState {
    errorDetails?: string;
    
}
export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = {};
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        const errData = JSON.stringify(error, Object.getOwnPropertyNames(error), 4);
        console.log(errData);
        return { errorDetails: errData };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
        console.error('An error occurred: ' + error, errorInfo);
    }

    render(): JSX.Element {
        if (this.state.errorDetails) {
            // Remove any modal backdrops in case the exception occurred in a modal
            document.querySelectorAll('div.modal-backdrop.fade.show').forEach(e => e.parentNode.removeChild(e));
            return (
                <div className="container">
                    <div className="card mt-3">
                        <div className="card-header">
                            An Error Occurred
                        </div>
                        <div className="card-body">
                            <p>An unrecoverable error occurred while attempting to render this page.
                                Please report this as an issue on <a href="https://github.com/ecnepsnai/otto/issues/new/choose" target="_blank" rel="noreferrer">Github</a> and include the following information:</p>
                            <pre>{this.state.errorDetails}</pre>
                        </div>
                    </div>
                </div>
            );
        }

        return (<React.Fragment>{this.props.children}</React.Fragment>);
    }
}
