import React from 'react';
import {observer} from 'mobx-react-lite';
import ScrollViewport, {IScrollableProps} from './ScrollViewport';
import {IScrollHook} from './scrollHooks';

export interface IProps<T> {
    className?: string;
    items: T[];
    renderItem: (item: T, index: number, selected: boolean) => React.ReactElement | null;
    selection: {
        row: number;
        onKeyDown: React.KeyboardEventHandler<HTMLElement>;
        onMouseDown: React.MouseEventHandler<HTMLElement>;
        scroll: Pick<IScrollHook<HTMLDivElement>, 'listRef' | 'rowHeight' | 'startRow' | 'endRow' | 'onScroll'>;
    };
    children?: React.ReactNode;
}

const ListViewport = observer(<T extends {id: number | string}>(props: IProps<T>) => {
    const {items, className, renderItem, selection, children} = props;
    const {startRow, rowHeight, onScroll, endRow, listRef} = selection.scroll;
    const leadingHeight = startRow * rowHeight;
    const height = items.length * rowHeight;
    return (
        <ScrollViewport onScroll={onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            {({scrollHeight}: IScrollableProps) => {
                if (scrollHeight === 0) return null;
                return (
                    <div ref={listRef} className={className} style={{height}}>
                        {leadingHeight > 0 ? <div style={{height: leadingHeight}} /> : null}
                        {items.slice(startRow, endRow(scrollHeight) + 1).map((row, index) =>
                            <React.Fragment key={row.id}>
                                {renderItem(row, startRow + index, startRow + index === selection.row)}
                            </React.Fragment>
                        )}
                        {children}
                    </div>
                );
            }}
        </ScrollViewport>
    );
});

export default ListViewport;
