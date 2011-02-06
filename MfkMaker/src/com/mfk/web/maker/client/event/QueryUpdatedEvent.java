package com.mfk.web.maker.client.event;

import com.google.gwt.event.shared.EventHandler;
import com.google.gwt.event.shared.GwtEvent;

public class QueryUpdatedEvent extends GwtEvent<QueryUpdatedEvent.Handler> {

  public interface Handler extends EventHandler {
    public void handleNewQuery(QueryUpdatedEvent event);
  }

  public static Type<Handler> TYPE = new Type<Handler>();

  public final String query;

  public QueryUpdatedEvent(String query) {
    this.query = query;
  }
  
  @Override
  public GwtEvent.Type<Handler> getAssociatedType() {
    return TYPE;
  }

  @Override
  protected void dispatch(Handler handler) {
    handler.handleNewQuery(this);
  }
}
