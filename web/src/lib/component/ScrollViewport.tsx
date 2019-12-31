import React from 'react';

type ResizeListener = (event: UIEvent) => void;

const resizeListeners: ResizeListener[] = [];
window.onresize = (event: UIEvent) => {
    resizeListeners.forEach(listener => {
        listener(event);
    });
};

interface IProps {
    className?: string;
    children: React.ReactNode;
    onScroll?: (event: React.UIEvent<HTMLElement>) => void;
}

export interface IScrollableProps {
    scrollHeight: number;
}

const ScrollViewport: React.FC<IProps> = ({children, className = 'scroll-container', onScroll}) => {
    const ref: React.MutableRefObject<HTMLDivElement> = React.useRef(null);
    const [scrollHeight, setScrollHeight] = React.useState(0);
    const onResize = React.useCallback(() => setScrollHeight(ref.current ? ref.current.clientHeight : 0), [ref.current]);
    React.useEffect(() => {
        if (ref.current) setScrollHeight(ref.current.clientHeight);
    }, [ref.current]);
    React.useEffect(() => {
        resizeListeners.push(onResize);
        return () => {
            const index = resizeListeners.indexOf(onResize);
            resizeListeners.splice(index, 1);
        };
    }, [onResize]);
    return (
        <div ref={ref} className={className} onScroll={onScroll}>
            {typeof children === 'function' ? children({scrollHeight}) : children}
        </div>);
};

export default ScrollViewport;