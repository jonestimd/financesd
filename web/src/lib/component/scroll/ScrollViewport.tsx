import React from 'react';

type ResizeListener = (event: UIEvent) => void;

const resizeListeners: ResizeListener[] = [];
window.onresize = (event: UIEvent) => {
    resizeListeners.forEach(listener => {
        listener(event);
    });
};

interface IProps extends React.DOMAttributes<HTMLDivElement> {
    className?: string;
    tabIndex?: number;
}

export interface IScrollableProps {
    scrollHeight: number;
}

const ScrollViewport: React.FC<IProps> = ({children, className = 'scroll-container', tabIndex = 0, ...attrs}) => {
    const ref = React.useRef<HTMLDivElement>(null);
    const [scrollHeight, setScrollHeight] = React.useState(0);
    const onResize = React.useCallback(() => setScrollHeight(ref.current ? ref.current.clientHeight : 0), [ref]);
    React.useEffect(() => {
        if (ref.current) setScrollHeight(ref.current.clientHeight);
    }, [ref]);
    React.useEffect(() => {
        resizeListeners.push(onResize);
        return () => {
            const index = resizeListeners.indexOf(onResize);
            resizeListeners.splice(index, 1);
        };
    }, [onResize]);
    return (
        <div ref={ref} className={className} {...attrs} tabIndex={tabIndex}>
            {typeof children === 'function' ? children({scrollHeight}) : children}
        </div>);
};

export default ScrollViewport;
