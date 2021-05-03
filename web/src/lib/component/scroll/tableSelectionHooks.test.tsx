import React from 'react';
import {shallow} from 'enzyme';
import {useSelection, ISelectionOptions} from './tableSelectionHooks';
import {createDiv} from 'src/test/htmlUtils';
import {mockScrollHook} from 'src/test/mockHooks';

const hook: Partial<ReturnType<typeof useSelection>> = {};

function setup(options: ISelectionOptions) {
    const TestComponent: React.FC = () => {
        Object.assign(hook, useSelection(options));
        return null;
    };
    const component = shallow(<TestComponent />);
    return component;
}
const document = new Document();

describe('tableSelectionHooks', () => {
    const rows = 1245;
    const columns = 4;
    const headerSelector = '.header';
    const rowSelector = '.row';
    const options = {rows, columns, rowSelector};
    const rowHeight = 10;

    describe('useSelection', () => {
        it('defaults to first row and column', () => {
            setup(options);

            expect(hook).toEqual(expect.objectContaining({row: 0, column: 0}));
        });
        it('defaults to scroll startRow', () => {
            mockScrollHook({startRow: 9});

            setup({...options});

            expect(hook).toEqual(expect.objectContaining({row: 9, column: 0}));
        });
        describe('setCell', () => {
            it('sets row and column', () => {
                mockScrollHook({startRow: 9});
                setup({...options});

                hook.setCell!(3, 4);

                expect(hook.row).toEqual(3);
                expect(hook.column).toEqual(4);
            });
        });
        describe('onKeyDown', () => {
            const container = createDiv();
            const row = createDiv();
            const header = createDiv();
            const children: Record<string, Element> = {
                [rowSelector]: row,
                [headerSelector]: header,
            };
            const containerHeight = 200;
            const headerHeight = 15;
            const pageSize = containerHeight / rowHeight;

            const mockKeyEvent = (key: string, ctrlKey: boolean, shiftKey = false) => ({
                currentTarget: container, key, ctrlKey, shiftKey,
                stopPropagation: jest.fn(),
                preventDefault: jest.fn(),
            } as unknown as React.KeyboardEvent<HTMLElement>);

            beforeEach(() => {
                jest.spyOn(container, 'getBoundingClientRect').mockReturnValue({height: containerHeight} as DOMRect);
                jest.spyOn(container, 'clientHeight', 'get').mockReturnValue(containerHeight);
                jest.spyOn(container, 'querySelector').mockImplementation((selector: keyof typeof children) => children[selector]);
                jest.spyOn(row, 'getBoundingClientRect').mockReturnValue({height: rowHeight} as DOMRect);
                jest.spyOn(header, 'getBoundingClientRect').mockReturnValue({height: headerHeight} as DOMRect);
            });
            it('stops event propagation', () => {
                setup(options);
                const event = mockKeyEvent('', false);

                hook.onKeyDown!(event);

                expect(event.stopPropagation).toBeCalledTimes(1);
                expect(event.preventDefault).not.toBeCalled();
                expect(container.getBoundingClientRect).toBeCalled();
                expect(container.querySelector).toBeCalledWith(rowSelector);
                expect(row.getBoundingClientRect).toBeCalled();
            });
            it('scrolls up to selected row', () => {
                const startRow = pageSize * 2;
                mockScrollHook({startRow});
                jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((startRow - pageSize / 2) * rowHeight);
                setup({...options, headerSelector});
                const event = mockKeyEvent('PageUp', false);

                hook.onKeyDown!(event);

                expect(container.scrollTo).toBeCalledWith({top: (startRow - pageSize) * rowHeight});
            });
            it('scrolls down to selected row', () => {
                const startRow = pageSize * 2;
                mockScrollHook({startRow});
                jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((startRow - pageSize / 2) * rowHeight);
                setup({...options, headerSelector});
                const event = mockKeyEvent('PageDown', false);

                hook.onKeyDown!(event);

                const finalRow = startRow + pageSize;
                expect(container.scrollTo).toBeCalledWith({top: finalRow * rowHeight + headerHeight - containerHeight + rowHeight});
            });
            const rowTests = [
                {name: 'selects previous row on up arrow', key: 'ArrowUp', ctrl: false, initialRow: 9, expectedRow: 8},
                {name: 'remains at top on up arrow', key: 'ArrowUp', ctrl: false, initialRow: 0, expectedRow: 0},
                {name: 'selects next row on down arrow', key: 'ArrowDown', ctrl: false, initialRow: 9, expectedRow: 10},
                {name: 'remains at bottom on down arrow', key: 'ArrowDown', ctrl: false, initialRow: rows - 1, expectedRow: rows - 1},
                {name: 'selects previous page on page up', key: 'PageUp', ctrl: false, initialRow: pageSize * 2, expectedRow: pageSize},
                {name: 'selects first row on page up', key: 'PageUp', ctrl: false, initialRow: pageSize * .5, expectedRow: 0},
                {name: 'selects next page on page down', key: 'PageDown', ctrl: false, initialRow: pageSize * 2, expectedRow: pageSize * 3},
                {name: 'selects last row on page down', key: 'PageDown', ctrl: false, initialRow: rows - pageSize * .5, expectedRow: rows - 1},
                {name: 'selects first row on ctrl+home', key: 'Home', ctrl: true, initialRow: rows - 1, expectedRow: 0},
                {name: 'remains on current row on home', key: 'Home', ctrl: false, initialRow: 9, expectedRow: 9},
                {name: 'selects last row on ctrl+end', key: 'End', ctrl: true, initialRow: 0, expectedRow: rows - 1},
                {name: 'remains on current row on end', key: 'End', ctrl: false, initialRow: 9, expectedRow: 9},
            ];
            rowTests.forEach(({name, initialRow, key, ctrl, expectedRow}) =>
                it(name, () => {
                    mockScrollHook({startRow: initialRow});
                    setup({...options});
                    const event = mockKeyEvent(key, ctrl);

                    hook.onKeyDown!(event);

                    expect(hook.row).toEqual(expectedRow);
                    expect(event.preventDefault).toBeCalled();
                })
            );

            const columnTests = [
                {name: 'selects previous column on left arrow', key: 'ArrowLeft', ctrl: false,
                    initial: {column: 2}, expected: {column: 1}},
                {name: 'wraps to last column on left arrow', key: 'ArrowLeft', ctrl: false,
                    initial: {column: 0}, expected: {column: columns - 1}},
                {name: 'selects next column on right arrow', key: 'ArrowRight', ctrl: false,
                    initial: {column: 0}, expected: {column: 1}},
                {name: 'wraps to first column on right arrow', key: 'ArrowRight', ctrl: false,
                    initial: {column: columns - 1}, expected: {column: 0}},
                {name: 'selects first column on home', key: 'Home', ctrl: false,
                    initial: {column: columns - 1}, expected: {column: 0}},
                {name: 'remains at current column on ctrl+home', key: 'Home', ctrl: true,
                    initial: {column: 1}, expected: {column: 1}},
                {name: 'selects last column on end', key: 'End', ctrl: false,
                    initial: {column: 0}, expected: {column: columns - 1}},
                {name: 'remains at current column on ctrl+end', key: 'End', ctrl: true,
                    initial: {column: 1}, expected: {column: 1}},
                {name: 'selects next column on tab', key: 'Tab', ctrl: false,
                    initial: {column: 0}, expected: {column: 1}},
                {name: 'wraps to next row on tab', key: 'Tab', ctrl: false,
                    initial: {column: columns - 1, row: 0}, expected: {column: 0, row: 1}},
                {name: 'remains on last cell row on tab', key: 'Tab', ctrl: false,
                    initial: {column: columns - 1, row: rows - 1}, expected: {column: columns - 1, row: rows - 1}},
                {name: 'selects previous column on shift tab', key: 'Tab', ctrl: false, shift: true,
                    initial: {column: 2}, expected: {column: 1}},
                {name: 'wraps to previous row on shift tab', key: 'Tab', ctrl: false, shift: true,
                    initial: {column: 0, row: 1}, expected: {column: columns - 1, row: 0}},
                {name: 'remains on first cell on shift tab', key: 'Tab', ctrl: false, shift: true,
                    initial: {column: 0, row: 0}, expected: {column: 0, row: 0}},
            ];
            columnTests.forEach(({name, initial, key, ctrl, shift, expected}) =>
                it(name, () => {
                    mockScrollHook({startRow: initial.row});
                    setup({...options});
                    for (let index = 0; index < initial.column; index++) hook.onKeyDown!(mockKeyEvent('ArrowRight', false));
                    const event = mockKeyEvent(key, ctrl, shift);

                    hook.onKeyDown!(event);

                    expect(hook.column).toEqual(expected.column);
                    if (expected.row) expect(hook.row).toEqual(expected.row);
                })
            );
        });
        describe('onMouseDown', () => {
            const container = document.createElement('div');
            const rows = new Array(50).fill(null).map(() => container.appendChild(document.createElement('div')));
            const target = document.createElement('span');
            rows[0].appendChild(target);

            const mockMouseEvent = (target: HTMLElement | null, clientY = -1, button = 0) => ({
                target,
                currentTarget: container,
                clientY,
                button,
                stopPropagation: jest.fn(),
                preventDefault: jest.fn(),
            } as unknown as React.MouseEvent<HTMLElement>);

            describe('table', () => {
                const cells = new Array(columns).fill(null).map(() => document.createElement('td'));
                const tr = document.createElement('tr');
                tr.setAttribute('class', 'row');
                cells.forEach((cell) => tr.appendChild(cell));
                const rowIndex = Math.floor(Math.random() * 100);

                beforeAll(() => {
                    cells.forEach((cell, index) => {
                        Object.defineProperty(cell, 'cellIndex', {value: index});
                        tr.appendChild(cell);
                    });
                    Object.defineProperty(tr, 'sectionRowIndex', {value: rowIndex});
                });
                it('selects row and column', () => {
                    const column = columns - 2;
                    cells.forEach((cell) => cell.colSpan = 1);
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(cells[column]));

                    expect(hook.row).toEqual(rowIndex);
                    expect(hook.column).toEqual(column);
                });
                it('uses colSpan to calculate column index', () => {
                    Object.defineProperty(cells[1], 'colSpan', {value: 2});
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(cells[2]));

                    expect(hook.row).toEqual(rowIndex);
                    expect(hook.column).toEqual(3);
                });
            });
            describe('non-table', () => {
                const boundingRect = (index: number) => ({bottom: (index + 1) * rowHeight} as DOMRect);

                beforeAll(() => rows.forEach((row) => row.setAttribute('class', rowSelector.substring(1))));
                beforeEach(() => {
                    rows.forEach((row, index) => jest.spyOn(row, 'getBoundingClientRect').mockReturnValue(boundingRect(index)));
                });
                it('selects row at click', () => {
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(target, rowHeight * 4.1));

                    expect(hook.row).toEqual(4);
                });
                it('ignores click past end row', () => {
                    const startRow = 5;
                    mockScrollHook({startRow});
                    setup({...options});

                    hook.onMouseDown!(mockMouseEvent(null, rowHeight * (rows.length + 1)));

                    expect(hook.row).toEqual(startRow);
                });
            });
        });
    });
});
