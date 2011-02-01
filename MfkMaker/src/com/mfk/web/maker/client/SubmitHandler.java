package com.mfk.web.maker.client;


import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.http.client.Request;
import com.google.gwt.http.client.RequestBuilder;
import com.google.gwt.http.client.RequestCallback;
import com.google.gwt.http.client.RequestException;
import com.google.gwt.http.client.Response;
import com.google.gwt.http.client.URL;
import com.google.gwt.http.client.UrlBuilder;
import com.google.gwt.user.client.Window;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.RootPanel;


// TODO(mjkelly): Do client-side validation here.
public class SubmitHandler implements ClickHandler {
  public void onClick(ClickEvent event) {
    MfkPanel[] p = { MfkMaker.items.get(0), MfkMaker.items.get(1),
        MfkMaker.items.get(2) };

    // Do some simple validation. This just to prevent mistakes. It
    // obviously does not replace real validation on the server side.
    if (p[0].similarTo(p[1]) || p[1].similarTo(p[2]) ||
        p[2].similarTo(p[0])) {
      MfkMaker.showDialog("Error Creating MFK",
          "Each item must have a unique title and image.");
      return;
    }
    for (int i = 0; i < 3; i++) {
      System.out.println("Checking URLs: " + p[i].image.getUrl() + " vs " + MfkMaker.defaultURLs[i]);
      if (p[i].image.getUrl().equals(MfkMaker.defaultURLs[i])) {
        MfkMaker.showDialog("Error Creating MFK",
            "You must change all the items from their defaults.");
        return;
      }
    }

    String url = "/rpc/create/";
    RequestBuilder builder = new RequestBuilder(RequestBuilder.POST, url);
    builder.setHeader("Content-Type", "application/x-www-form-urlencoded");
    StringBuffer reqData = new StringBuffer();

    URL.encodeQueryString(url);
    reqData.append("n1=").append(URL.encodeQueryString(p[0].title));
    reqData.append("&n2=").append(URL.encodeQueryString(p[1].title));
    reqData.append("&n3=").append(URL.encodeQueryString(p[2].title));
    reqData.append("&u1=").append(
        URL.encodeQueryString(p[0].image.getUrl()));
    reqData.append("&u2=").append(
        URL.encodeQueryString(p[1].image.getUrl()));
    reqData.append("&u3=").append(
        URL.encodeQueryString(p[2].image.getUrl()));
    reqData.append("&q1=").append(
        URL.encodeQueryString(p[0].image.getQuery()));
    reqData.append("&q2=").append(
        URL.encodeQueryString(p[1].image.getQuery()));
    reqData.append("&q3=").append(
        URL.encodeQueryString(p[2].image.getQuery()));
    System.out.println("request data = " + reqData);

    try {
      builder.sendRequest(reqData.toString(), new RequestCallback() {
        public void onError(Request request, Throwable exception) {
          System.out.println("Error creating new Triple");
        }

        public void onResponseReceived(Request request,
            Response response) {
          if (response.getStatusCode() == 200) {
            String[] responseParts = response.getText().split(":",
                2);
            if (responseParts[0].equals("ok")) {
              System.out.println("Successful creation request: "
                  + response.getText());

              // clear the existing page in a haphazard way
              MfkMaker.itemPanel.clear();
              RootPanel.get("submit").clear();
              RootPanel.get("instructions").setVisible(false);

              UrlBuilder builder = Window.Location.createUrlBuilder();
              builder.setPath("/");
              builder.setParameter("id", responseParts[1]);
              String url = builder.buildString();
              MfkMaker.itemPanel.add(
                  new HTML("<h2>Great Success!</h2>"
                  + "You successfully created a new "
                  + "MFK item. Other users can vote "
                  + "on it here:"
                  + "<p><a href='" + url + "'>"
                  + url + "</a></p>"
                  + "<p>Thanks for contributing! "
                  + "<a href='/make'>Make another?</a></p>"));


            } else {
              System.out.println("Error: " + responseParts[1]);
              MfkMaker.showDialog("Error Creating Your MFK",
                  "The server says: " + responseParts[1]);
            }
          } else {
            System.out.println(
                "Server-side error creating new MFK: "
                + "Response code: "
                + response.getStatusCode()
                + "; response text: "
                + response.getText());
            MfkMaker.showDialog("Server Error",
                "Response code "
                + response.getStatusCode() + ": "
                + response.getStatusText());
          }
        }
      });
    } catch (RequestException e) {
      System.out.println("Error sending vote: " + e);
    }

  }
}
