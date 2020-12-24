import React from 'react';
import {observer} from 'mobx-react-lite';
import ScrollViewport, {IScrollableProps} from './ScrollViewport';
import {useScroll} from './scrollHooks';
// import {useSelection} from './table/selection';

interface IProps<T> {
    className?: string;
    defaultRowHeight?: number;
    prototypeSelector?: string;
    items: T[];
    renderItem: (item: T, index: number, selected: boolean) => React.ReactElement | null;
    children?: React.ReactNode;
}

const ListViewPort = observer(<T extends {id: number | string}>(props: IProps<T>) => {
    const {items, className, defaultRowHeight, prototypeSelector, renderItem, children} = props;
    const scroll = useScroll<HTMLDivElement>({defaultRowHeight, prototypeSelector});
    const leadingHeight = scroll.startRow * scroll.rowHeight;
    const height = items.length * scroll.rowHeight;
    // const selection = useSelection(scroll.startRow, items.length, 1, Math.max(0, scroll.startRow - 1));
    return (
        // <ScrollViewport onScroll={scroll.onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
        <ScrollViewport onScroll={scroll.onScroll}>
            {({scrollHeight}: IScrollableProps) => {
                if (scrollHeight === 0) return null;
                const endRow = scroll.startRow + Math.ceil(scrollHeight / scroll.rowHeight) * 3;
                return (
                    <div ref={scroll.listRef} className={className} style={{height}}>
                        {children}
                        {leadingHeight > 0 ? <div style={{height: leadingHeight}} /> : null}
                        {items.slice(scroll.startRow, endRow + 1).map((row, index) =>
                            <React.Fragment key={row.id}>
                                {/* {renderItem(row, scroll.startRow + index, scroll.startRow + index === selection.row)} */}
                                {renderItem(row, scroll.startRow + index, false)}
                            </React.Fragment>
                        )}
                    </div>
                );
            }}
        </ScrollViewport>
    );
});

export default ListViewPort;