FROM scratch
ADD  eventscollector /
ENV PORT=13000
EXPOSE 13000
CMD ["/eventscollector"]