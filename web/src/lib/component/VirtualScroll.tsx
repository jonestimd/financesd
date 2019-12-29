import React from 'react';
import classNames from 'classnames';

interface IProps {
    className?: string;
    children: React.ReactNode;
    itemCount: number;
    start: number;
    end: number;
    onScroll: (deltaY: number) => void;
}

const scrollSize = 15;
// window.onresize = (event: UIEvent) => console.log(event);

const VitrualScroll: React.FC<IProps> = ({children, itemCount, start, end, onScroll, className = 'scroll-container'}) => {
    const onWheel = React.useCallback(({deltaY}: React.WheelEvent) => onScroll(deltaY), [onScroll]);
    const scrollbarRef: React.MutableRefObject<HTMLDivElement> = React.useRef(null);
    const viewHeight = scrollbarRef.current ? scrollbarRef.current.clientHeight - 2 * scrollSize : 0;
    const height = itemCount ? Math.max(scrollSize, viewHeight * (end - start) / itemCount) : viewHeight;
    const top = itemCount ? Math.min(viewHeight - height, start * viewHeight / itemCount) : 0;
    return (
        <div className={classNames(className, 'hide-scroll')} onWheel={onWheel}>
            <div className='scrollable'>{children}</div>
            <div ref={scrollbarRef} className='scroll-bar'>
                <div className='up arrow' />
                <div className='scroll-handle' style={{height, top}} />
                <div className='down arrow' />
            </div>
        </div>);
};

export default VitrualScroll;