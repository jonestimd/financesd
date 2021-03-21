import {reaction} from 'mobx';

export function testAction(get: () => void, actionToTest: () => void) {
    const listen = jest.fn();
    const dispose = reaction(get, listen);
    actionToTest();
    dispose();
    expect(listen).toBeCalledTimes(1);
}
