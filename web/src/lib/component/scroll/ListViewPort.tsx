import React from 'react';
import {observer} from 'mobx-react-lite';
import ScrollViewport, {IScrollableProps} from './ScrollViewport';
import {useScroll} from './scrollHooks';
import {useSelection} from './selectionHooks';

interface IProps<T> {
    className?: string;
    defaultRowHeight?: number;
    rowSelector: string;
    prototypeSelector?: string;
    items: T[];
    renderItem: (item: T, index: number, selected: boolean) => React.ReactElement | null;
    children?: React.ReactNode;
}

const ListViewPort = observer(<T extends {id: number | string}>(props: IProps<T>) => {
    const {items, className, defaultRowHeight, rowSelector, prototypeSelector, renderItem, children} = props;
    const scroll = useScroll<HTMLDivElement>({defaultRowHeight, prototypeSelector});
    const leadingHeight = scroll.startRow * scroll.rowHeight;
    const height = items.length * scroll.rowHeight;
    const selection = useSelection({initialRow: scroll.startRow, rows: items.length, rowOffset: scroll.startRow, rowSelector});
    return (
        <ScrollViewport onScroll={scroll.onScroll} onKeyDown={selection.onKeyDown} onMouseDown={selection.onMouseDown}>
            {({scrollHeight}: IScrollableProps) => {
                if (scrollHeight === 0) return null;
                const endRow = scroll.endRow(scrollHeight);
                return (
                    <div ref={scroll.listRef} className={className} style={{height}}>
                        {leadingHeight > 0 ? <div style={{height: leadingHeight}} /> : null}
                        {items.slice(scroll.startRow, endRow + 1).map((row, index) =>
                            <React.Fragment key={row.id}>
                                {renderItem(row, scroll.startRow + index, scroll.startRow + index === selection.row)}
                            </React.Fragment>
                        )}
                        {children}
                    </div>
                );
            }}
        </ScrollViewport>
    );
});

export default ListViewPort;
