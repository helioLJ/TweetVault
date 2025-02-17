import mitt from 'mitt';

type Events = {
  tagUpdated: void;
};

export const eventBus = mitt<Events>(); 