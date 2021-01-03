import React from 'react';
import {shallow} from 'enzyme';
import ListViewport from './ListViewport';
import ScrollViewport from './ScrollViewport';
import * as scrollHooks from './scrollHooks';
import * as selectionHooks from './selectionHooks';

type ScrollHook = Pick<ReturnType<typeof scrollHooks['useScroll']>, 'startRow' | 'rowHeight' | 'headerHeight'>;

function mockHooks(overrides: Partial<ScrollHook> = {}) {
    const scroll = {
        startRow: 0,
        rowHeight: 0,
        headerHeight: 0,
        listRef: {current: null} as React.MutableRefObject<HTMLElement>,
        endRow: jest.fn<number, [number]>(),
        onScroll: jest.fn<void, [React.UIEvent<HTMLElement>]>(),
        ...overrides,
    };
    jest.spyOn(scrollHooks, 'useScroll').mockReturnValue(scroll);
    const selection = {
        row: 0,
        column: 0,
        onKeyDown: jest.fn(),
        onMouseDown: jest.fn(),
    };
    jest.spyOn(selectionHooks, 'useSelection').mockReturnValue(selection);
    return {scroll, selection};
}

describe('ListViewport', () => {
    const rowHeight = 10;
    const props = {
        className: 'container',
        rowSelector: '.row',
        renderItem: (row: {id: number}) => <div className='row'>{row.id}</div>,
        items: new Array(10).fill(null).map((_, index) => ({id: index + 1})),
    };

    it('returns null for 0 height', () => {
        const {scroll, selection} = mockHooks({rowHeight: 10});

        const viewport = shallow(<ListViewport {...props} />).find(ScrollViewport);

        expect(viewport).toHaveProp('onScroll', scroll.onScroll);
        expect(viewport).toHaveProp('onKeyDown', selection.onKeyDown);
        expect(viewport).toHaveProp('onMouseDown', selection.onMouseDown);
        expect(viewport.prop('children')({scrollHeight: 0})).toBeNull();
    });
    it('does not display leading filler if startRow is 0', () => {
        const {scroll} = mockHooks({rowHeight});
        const endRow = 6;
        scroll.endRow.mockReturnValue(endRow);
        const component = shallow(<ListViewport {...props} />);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        const container = result.find('.container');
        const children = container.prop<React.ReactNode[]>('children');
        expect(children[0]).toBeNull();
        expect(children[1]).toHaveLength(endRow + 1);
    });
    it('displays items from scroll.startRow to scroll.endRow', () => {
        const startRow = 2;
        const endRow = 6;
        const {scroll} = mockHooks({rowHeight, startRow});
        scroll.endRow.mockReturnValue(endRow);
        const component = shallow(<ListViewport {...props} />);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        expect(result.find('.container')).toHaveStyle({height: props.items.length * rowHeight});
        expect(result.find('.row').map((r) => r.text()))
            .toEqual(props.items.slice(startRow, endRow + 1).map((item) => `${item.id}`));
        expect(result.find('.container').childAt(0)).toHaveStyle({height: startRow * rowHeight});
    });
    it('displays children after rows', () => {
        const startRow = 2;
        const endRow = 6;
        const {scroll} = mockHooks({rowHeight, startRow});
        scroll.endRow.mockReturnValue(endRow);
        const component = shallow(<ListViewport {...props}><div className='prototype'>a sample row</div></ListViewport>);

        const result = shallow(component.find(ScrollViewport).prop('children')({scrollHeight: 100}) as React.ReactElement);

        expect(result.find('.row + .prototype')).toHaveText('a sample row');
    });
});
