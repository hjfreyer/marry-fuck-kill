
package com.mfk.web.maker.client.event;

import java.util.List;

import com.google.gwt.event.shared.EventHandler;
import com.google.gwt.event.shared.GwtEvent;

public class ImageResultsAvailableEvent extends GwtEvent<ImageResultsAvailableEvent.Handler> {

  public interface Handler extends EventHandler {
    public void handleImageResults(ImageResultsAvailableEvent event);
  }

  public static Type<Handler> TYPE = new Type<Handler>();

  public final boolean newResultSet;  
  public final String query;
  public final List<String> resultUrls;
  
  public ImageResultsAvailableEvent(boolean newQuery, String query,
      List<String> resultUrls) {
    this.newResultSet = newQuery;
    this.query = query;
    this.resultUrls = resultUrls;
  }

  @Override
  public GwtEvent.Type<Handler> getAssociatedType() {
    return TYPE;
  }

  @Override
  protected void dispatch(Handler handler) {
    handler.handleImageResults(this);
  }
}
