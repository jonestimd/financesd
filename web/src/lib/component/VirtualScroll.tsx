import React from 'react';
import classNames from 'classnames';

type ResizeListener = (event: UIEvent) => void;

interface IProps {
    className?: string;
    children: React.ReactNode;
    itemCount: number;
    start: number;
    end: number;
    onScroll: (deltaY: number) => void;
    onResize: ResizeListener;
}

const scrollSize = 15;

const resizeListeners: ResizeListener[] = [];
window.onresize = (event: UIEvent) => {
    resizeListeners.forEach(listener => {
        listener(event);
    });
};

const VitrualScroll: React.FC<IProps> = ({children, itemCount, start, end, onScroll, onResize, className = 'scroll-container'}) => {
    React.useEffect(() => {
        resizeListeners.push(onResize);
        return () => {
            const index = resizeListeners.indexOf(onResize);
            resizeListeners.splice(index, 1);
        };
    }, [onResize]);
    const onWheel = React.useCallback(({deltaY}: React.WheelEvent) => onScroll(deltaY), [onScroll]);
    const onKeyDown = React.useCallback(({currentTarget: target, key, shiftKey, altKey, ctrlKey, metaKey}: React.KeyboardEvent) => {
        if (!shiftKey && !altKey && !metaKey) {
            if (ctrlKey) {
                if (key === 'Home') onScroll(-Infinity);
                else if (key === 'End') onScroll(Infinity);
            }
            else {
                if (key === 'PageUp') onScroll(-target.clientHeight / 2);
                else if (key === 'PageDown') onScroll(target.clientHeight / 2);
            }
        }
    }, [onScroll]);
    const scrollbarRef: React.MutableRefObject<HTMLDivElement> = React.useRef(null);
    const viewHeight = scrollbarRef.current ? scrollbarRef.current.clientHeight - 2 * scrollSize : 0;
    const handleHeight = itemCount ? viewHeight * (end - start) / itemCount : viewHeight;
    const top = itemCount ? Math.min(viewHeight - handleHeight, start * viewHeight / itemCount) : 0;
    return (
        <div className={classNames(className, 'hide-scroll')} onWheel={onWheel}>
            <div className='scrollable' onKeyDown={onKeyDown} tabIndex={0}>{children}</div>
            <div ref={scrollbarRef} className='scroll-bar'>
                <div className='up arrow' />
                <div className='scroll-handle' style={{height: handleHeight, top}} />
                <div className='down arrow' />
            </div>
        </div>);
};

export default VitrualScroll;