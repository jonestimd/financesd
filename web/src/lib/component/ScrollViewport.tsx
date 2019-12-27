import React from 'react';
import classNames from 'classnames';

interface IProps {
    className?: string;
    children: React.ReactNode;
    scroll?: {
        itemCount: number;
        start: number;
        end: number;
        onWheel: (event: React.WheelEvent) => void;
    };
}

const scrollSize = 15;
// window.onresize = (event: UIEvent) => console.log(event);

const ScrollViewport: React.FC<IProps> = ({children, scroll: control, className = 'scroll-container'}) => {
    if (control) {
        const viewRef: React.MutableRefObject<HTMLDivElement> = React.useRef(null);
        const {onWheel, itemCount, start, end} = control || {};
        const viewHeight = viewRef.current ? viewRef.current.clientHeight - 2 * scrollSize : 0;
        const height = itemCount ? viewHeight / itemCount * (end - start) : viewHeight;
        const top = itemCount ? start * viewHeight / itemCount : 0;
        return (
            <div className={classNames(className, 'hide-scroll')} onWheel={onWheel}>
                <div className='scrollable'>{children}</div>
                <div ref={viewRef} className='scroll-bar'>
                    <div className='up arrow' />
                    <div className='scroll-handle' style={{height, top}}/>
                    <div className='down arrow' />
                </div>
            </div>);
    }
    else return <div className={className}>{children}</div>;
};

export default ScrollViewport;