import React from 'react';
import {shallow} from 'enzyme';
import {useSelection, ISelectionOptions, ICell} from './listSelectionHooks';
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

describe('listSelectionHooks', () => {
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
                {name: 'selects last row on ctrl+end', key: 'End', ctrl: true, initialRow: 0, expectedRow: rows - 1},
            ];
            rowTests.forEach(({name, initialRow, key, ctrl, expectedRow}) =>
                it(name, () => {
                    mockScrollHook({startRow: initialRow});
                    const getColumn = jest.fn<number, [ICell]>().mockImplementation(({column}) => column+1);
                    setup({...options, getColumn});
                    const event = mockKeyEvent(key, ctrl);

                    hook.onKeyDown!(event);

                    expect(hook.row).toEqual(expectedRow);
                    expect(event.stopPropagation).toBeCalled();
                    expect(event.preventDefault).toBeCalled();
                    expect(getColumn).toBeCalledWith({row: expectedRow, column: 0});
                    expect(hook.column).toEqual(1);
                })
            );
            const homeEndTests = [
                {name: 'remains on current row on home', key: 'Home', initialRow: 9},
                {name: 'remains on current row on end', key: 'End', initialRow: 9},
            ];
            homeEndTests.forEach(({name, initialRow, key}) =>
                it(name, () => {
                    mockScrollHook({startRow: initialRow});
                    setup({...options});
                    const event = mockKeyEvent(key, false);

                    hook.onKeyDown!(event);

                    expect(hook.row).toEqual(initialRow);
                    expect(event.stopPropagation).not.toBeCalled();
                    expect(event.preventDefault).not.toBeCalled();
                })
            );
        });
        describe('onMouseDown', () => {
            const container = document.createElement('div');
            const rows = new Array(50).fill(null).map(() => container.appendChild(document.createElement('div')));
            const target = document.createElement('span');
            rows[0].appendChild(target);

            const mockMouseEvent = (target: HTMLElement, clientY = -1, button = 0) => ({
                target,
                currentTarget: container,
                clientY,
                button,
                stopPropagation: jest.fn(),
                preventDefault: jest.fn(),
            } as unknown as React.MouseEvent<HTMLElement>);

            const boundingRect = (index: number) => ({bottom: (index + 1) * rowHeight} as DOMRect);

            beforeAll(() => rows.forEach((row) => row.setAttribute('class', rowSelector.substring(1))));
            beforeEach(() => {
                rows.forEach((row, index) => jest.spyOn(row, 'getBoundingClientRect').mockReturnValue(boundingRect(index)));
            });
            it('selects row at click', () => {
                const startRow = 5;
                mockScrollHook({startRow});
                const getColumn = jest.fn<number, [ICell]>().mockImplementation(({column}) => column+1);
                setup({...options, getColumn});

                hook.onMouseDown!(mockMouseEvent(target, rowHeight * 4.1));

                expect(hook.row).toEqual(4 + startRow);
                expect(hook.column).toEqual(1);
                expect(getColumn).toBeCalledWith({row: 4 + startRow, column: 0});
            });
            it('ignores click on current row', () => {
                const startRow = 5;
                mockScrollHook({startRow});
                const getColumn = jest.fn<number, [ICell]>().mockImplementation(({column}) => column+1);
                setup({...options, getColumn});

                hook.onMouseDown!(mockMouseEvent(target, rowHeight * 0.1));

                expect(hook.row).toEqual(startRow);
                expect(hook.column).toEqual(0);
                expect(getColumn).not.toBeCalled();
            });
            it('ignores click past end row', () => {
                const startRow = 5;
                mockScrollHook({startRow});
                setup({...options});

                hook.onMouseDown!(mockMouseEvent(target, rowHeight * (rows.length + 1)));

                expect(hook.row).toEqual(startRow);
            });
            it('ignores click of wrong button', () => {
                const startRow = 5;
                mockScrollHook({startRow});
                setup({...options});

                hook.onMouseDown!(mockMouseEvent(target, rowHeight * 2, 1));

                expect(hook.row).toEqual(startRow);
            });
            it('ignores click outside container', () => {
                const startRow = 5;
                mockScrollHook({startRow});
                setup({...options});

                hook.onMouseDown!(mockMouseEvent(document.createElement('span'), rowHeight * 2));

                expect(hook.row).toEqual(startRow);
            });
        });
    });
});
