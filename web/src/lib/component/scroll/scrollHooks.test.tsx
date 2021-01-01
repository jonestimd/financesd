import React from 'react';
import {shallow} from 'enzyme';
import {defaultOverscan, useScroll, IScrollOptions} from './scrollHooks';

const hook: Partial<ReturnType<typeof useScroll>> = {};

function setup(options: IScrollOptions) {
    const TestComponent: React.FC = () => {
        Object.assign(hook, useScroll(options));
        return null;
    };
    const component = shallow(<TestComponent />);
    return component;
}

describe('scrollHooks', () => {
    const defaultRowHeight = 24;
    const document = new HTMLDocument();
    const list = document.createElement('div');
    const prototype = document.createElement('span');

    describe('useScroll', () => {
        describe('startRow', () => {
            it('defaults to 0', () => {
                setup({});

                expect(hook.startRow).toEqual(0);
            });
        });
        describe('rowHeight', () => {
            it('returns default row height', () => {
                setup({defaultRowHeight: 42});

                expect(hook.rowHeight).toEqual(42);
            });
            it('calculates height from prototype', () => {
                const prototypeSelector = '.prototype';
                jest.spyOn(list, 'querySelector').mockReturnValue(prototype);
                jest.spyOn(prototype, 'getBoundingClientRect').mockReturnValue({height: 42} as DOMRect);
                const component = setup({prototypeSelector});

                hook.listRef.current = list;
                component.rerender({x: 'force rerender'});

                expect(hook.rowHeight).toEqual(42);
                expect(list.querySelector).toBeCalledWith(prototypeSelector);
            });
        });
        describe('headerHeight', () => {
            it('defaults to 0', () => {
                setup({});

                expect(hook.headerHeight).toEqual(0);
            });
            it('calculates height using header selector', () => {
                const headerSelector = '.header';
                jest.spyOn(list, 'querySelector').mockReturnValue(prototype);
                jest.spyOn(prototype, 'getBoundingClientRect').mockReturnValue({height: 42} as DOMRect);
                const component = setup({headerSelector});

                hook.listRef.current = list;
                component.rerender({x: 'force rerender'});

                expect(hook.headerHeight).toEqual(42);
                expect(list.querySelector).toBeCalledWith(headerSelector);
            });
        });
        describe('endRow', () => {
            it('calculates end row for overscan', () => {
                const visibleRows = 15.4;

                setup({defaultRowHeight});

                expect(hook.endRow(defaultRowHeight * visibleRows)).toEqual(Math.ceil(visibleRows * (1 + 2 * defaultOverscan)));
            });
        });
        describe('onScroll', () => {
            function testOnScroll(leadingRows: number, scrollTopRow: number, expectedStartRow: number) {
                const visibleRows = leadingRows / defaultOverscan;
                const component = setup({defaultRowHeight});
                jest.spyOn(list, 'clientHeight', 'get').mockReturnValue(defaultRowHeight * visibleRows);
                jest.spyOn(list, 'scrollTop', 'get').mockReturnValue(defaultRowHeight * (scrollTopRow + 0.4));

                hook.onScroll({currentTarget: list} as unknown as React.UIEvent<HTMLElement>);
                component.rerender({x: 'force rerender'});

                expect(hook.startRow).toEqual(expectedStartRow);
            }

            it('updates start row', () => {
                testOnScroll(2, 20, 18);
            });
            it('forces even start row', () => {
                testOnScroll(3, 20, 18);
            });
        });
    });
});
