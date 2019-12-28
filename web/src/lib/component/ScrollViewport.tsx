import React from 'react';

interface IProps {
    className?: string;
    children: React.ReactNode;
}

const ScrollViewport: React.FC<IProps> = ({children, className = 'scroll-container'}) => {
    return <div className={className}>{children}</div>;
};

export default ScrollViewport;