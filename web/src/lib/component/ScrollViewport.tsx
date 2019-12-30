import React from 'react';

interface IProps {
    className?: string;
    children: React.ReactNode;
    onScroll?: (event: React.UIEvent<HTMLElement>) => void;
}

const ScrollViewport: React.FC<IProps> = ({children, className = 'scroll-container', onScroll}) => {
    return <div className={className} onScroll={onScroll}>{children}</div>;
};

export default ScrollViewport;