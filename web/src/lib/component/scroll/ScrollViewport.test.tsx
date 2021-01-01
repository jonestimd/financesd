import React from 'react';
import {shallow} from 'enzyme';
import ScrollViewport, {IScrollableProps} from './ScrollViewport';

describe('ScrollViewport', () => {
    const document = new Document();
    const disposers: (() => void)[] = [];
    const capture = {
        height: -1,
        setHeight: ({scrollHeight}: IScrollableProps): React.ReactNode => {
            capture.height = scrollHeight;
            return null;
        },
    };

    beforeEach(() => {
        const effects: string[] = [];
        jest.spyOn(React, 'useEffect').mockImplementation((effect) => {
            if (!effects.includes(effect.toString())) {
                const disposer = effect();
                if (disposer) disposers.push(disposer);
                effects.push(effect.toString());
            }
        });
    });
    afterEach(() => {
        disposers.forEach((disposer) => disposer());
        disposers.splice(0, disposers.length);
    });
    it('defaults height to 0 if not mounted', () => {
        shallow(<ScrollViewport>{capture.setHeight}</ScrollViewport>);

        expect(capture.height).toEqual(0);
    });
    it('sets initial height', () => {
        const ref = {current: document.createElement('div')} as React.MutableRefObject<HTMLDivElement>;
        jest.spyOn(ref.current, 'clientHeight', 'get').mockReturnValue(100);
        jest.spyOn(React, 'useRef').mockReturnValue(ref);

        shallow(<ScrollViewport>{capture.setHeight}</ScrollViewport>);

        expect(capture.height).toEqual(100);
    });
    it('updates height when window is resized', () => {
        const ref = {current: document.createElement('div')} as React.MutableRefObject<HTMLDivElement>;
        jest.spyOn(ref.current, 'clientHeight', 'get').mockReturnValueOnce(100);
        jest.spyOn(ref.current, 'clientHeight', 'get').mockReturnValueOnce(150);
        jest.spyOn(React, 'useRef').mockReturnValue(ref);
        shallow(<ScrollViewport>{capture.setHeight}</ScrollViewport>);

        window.onresize(null);

        expect(capture.height).toEqual(150);
    });
});
